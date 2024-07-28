import yt_dlp
import json
from openai import OpenAI
import os
from dotenv import load_dotenv
load_dotenv()

PLAYLIST_LINK = "https://www.youtube.com/playlist?list=PLNOhvqcJZLWjXdfJK4wCfZVzxq0VyCpoh"
TRANSCRIPTION_CORRECTION_SYSTEM_PROMPT = """
You are a helpful assistant correcting video transcriptions that instruct the viewer on a task.
Your task is to correct any spelling discrepancies in the transcribed text. 
Exclude the mention of "Today's Mission" at the beginning.
"""

def download_video(video_url, name, output_path):
    output_folder = os.path.join(output_path, name)
    if not os.path.exists(output_folder):
    # If it doesn't exist, create the folder
        os.makedirs(output_folder)

    # Download audio as WAV
    ydl_opts_audio = {
        'format': 'bestaudio/best',
        'postprocessors': [{
            'key': 'FFmpegExtractAudio',
            'preferredcodec': 'wav',
            'preferredquality': '192',
        }],
        'outtmpl': os.path.join(output_folder, 'audio.%(ext)s'),
        'ignoreerrors': True,
    }

    with yt_dlp.YoutubeDL(ydl_opts_audio) as ydl:
        ydl.download([video_url])
    
    return os.path.join(output_folder, 'audio.wav')
   
def get_playlist_videos(playlist_url):
    ydl_opts = {
        'extract_flat': True,
        'ignore_errors': True,
    }

    playlist_data = {}

    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        playlist_info = ydl.extract_info(playlist_url, download=False)
        if 'entries' in playlist_info:
            for entry in playlist_info['entries']:
                if entry['duration'] != None and entry['duration'] < 60:
                    title = entry['title'].replace(':  ', ': ')
                    playlist_data[title] = {'title': title, 'url': entry['url']}
    return playlist_data

def transcribe(file_path):
    def correct_transcript(transcription):
        response = client.chat.completions.create(
            model="gpt-4o",
            messages=[
                {
                    "role": "system",
                    "content": TRANSCRIPTION_CORRECTION_SYSTEM_PROMPT
                },
                {
                    "role": "user",
                    "content": transcription
                }
            ]
        )   
        print(response.choices[0].message.content)
        return response.choices[0].message.content
    
    client = OpenAI()

    audio_file= open(file_path, "rb")
    transcription = client.audio.transcriptions.create(
        model="whisper-1", 
        file=audio_file,
        prompt="TODAY'S MISSION:"
    )
    os.remove(file_path)
    return correct_transcript(transcription.text)

def store_data(video_data_entry, output_path):
    output_folder = os.path.join(output_path, video_data_entry['title'])
    with open(os.path.join(output_folder, 'data.json'), 'w') as json_file:
        json.dump(video_data_entry, json_file, indent=2)

def download_playlist(playlist_url, output_path):
    video_data = get_playlist_videos(playlist_url)
    for i, video in enumerate(list(video_data.keys()), 1):
        print(f"Downloading video {i} of {len(video_data)}")
        audio_file_path = download_video(video_data[video]['url'], video, output_path)
        video_data[video]['text'] = transcribe(audio_file_path)
        video_data[video]['id'] = video_data[video]['url'].replace("https://www.youtube.com/watch?v=", "")
        store_data(video_data[video], output_path)


download_playlist(PLAYLIST_LINK, "media")