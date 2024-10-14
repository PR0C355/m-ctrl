- [Mission Control aka `m-ctrl`](#mission-control-aka-m-ctrl)
  - [Inspiration](#inspiration)
  - [How It Works](#how-it-works)
    - [Transcription](#transcription)
    - [The API](#the-api)
    - [The React Application](#the-react-application)
  - [Deployment](#deployment)
  - [Further Improvements](#further-improvements)

# Mission Control aka `m-ctrl`

Mission Control is a full-stack application driven by 3 main objectives:

1. Learn to use `go` to develop a simple API
2. Develop a very simple React.js application
3. Create a fun way to add simple personal challenges/objectives to everyday life

## Inspiration
The primary inspiration behind this application was to dip my toe into the world of both React development, as well as leveraging a statically typed, compiled language like `go` to develop a RESTful API.

The idea, however, came from a series of videos from a channel called ["FREAKBAiT" on YouTube](https://www.youtube.com/@FREAKBAiT). This series of videos is called "TODAY'S MISSION" wherein which a narrator usually instructs the viewer to engage in a certain behavior or action at some point in that day that may usually be considered weird or odd in standard social settings. 

This can range from things like "Superglue a coin to the ground and see how many people try to pick it up." to pranking your friends with something like "Send your friend a fake letter from a celebrity offering them a dream job."

I've found many of these videos funny or entertaining, but due to the nature of our fast-moving world, often forget whatever challenge is posed to me by these videos within minutes. Thus, I set out to create a platform where I can easily access one of these missions and remind myself of what the challenge is in order to complete it by the end of the day, similar to [Josh Wardle's](https://github.com/powerlanguage) [Wordle](https://www.nytimes.com/games/wordle/index.html). 

## How It Works
The application is split into three main parts:
 1. [Transcription](#transcription)
 2. [The API](#the-api)
 3. [The React Application](#the-react-application)


### Transcription
The transcription process is relatively simple, first the ["TODAY'S MISSION" playlist](https://www.youtube.com/playlist?list=PLNOhvqcJZLWjXdfJK4wCfZVzxq0VyCpoh) created by the FREAKBAiT channel on YouTube is scraped for it's contents, getting the YouTube link for each video, as well as the video's title.

Then, the audio for each video is saved locally. This is then fed into OpenAI's whisper API to transcribe the speech into text. 

After that, it is sent through GPT-4 with a custom transcription correction prompt, to increase the accuracy of the transcriptions. 

Finally, the original audio file, as well as the resulting transcription, video title, and video URL in a `JSON` file, are saved in folders named after each of the videos within a specified directory. 

### The API
This `go` application implements a simple API server for managing and retrieving mission data, using the Gin web framework for routing and handling HTTP requests.

Firstly, the service starts by loading the `JSON` files from a specific directory, that contain the transcriptions and video URLs. This means, however, that mission data is loaded into memory at startup and is not persisted or updated during runtime.

The intended usage is to provide both this and the python transcription script with shared filespace, so that the transcription script may store the results lcoally, and the API can load them up from the same location.

The API has a few different endpoints:

| Endpoint               | Purpose                                                     |
| ---------------------- | ----------------------------------------------------------- |
| `/missions/get/all`    | Retrieve all stored missions                                |
| `/missions/get/random` | Retrieve a random mission                                   |
| `/missions/get/unique` | Retrieve a unique mission based on the provided query token |
| `/missions/get/{id}`   | Retrieve a specific mission based on it's ID                |

At the time of writing, only the `/missions/get/unique` endpoint is used by the frontend application.

Most of these endpoints are relatively self-explanatory. However, `/missions/get/unique` is, well, **unique**. It functions by taking a session token that should be provided to the API within the query, which is then used to generate a has based on the day, which is then used as a seed for a random number generator, and the first number from the generator is then used to select from the list of existing missions.

The result is that any two requests with the same session token will return the same entry/mission, but if the two differ by so much as one character, then the result will likely be completely different. Additionally, because the hash is generated based on the date, the same session token will likely have differeint results on two different days. 

Thus, even if the token doesn't change, the result tomorrow will be different from that of today, but multiple requests within the same day will still return the same result.

### The React Application
The frontend is a very simple single-page React Progressive Web App (PWA). 

Upon access, the app loads up default Mission info, including a title, a transcription, and video. However, since these are never meant to be seen they are very clearly placeholders. Then it attempts to fetch a unique mission for the user. 

This is done by first checking the browser's local storage for a session token, and if one is missing, a new one is generated from a randomly generated token, and the date, creating a unique hash. This hash is then passed as the token in the query sent to [the API's](#the-api) `/missions/get/unique` endpoint, resulting in a response from the server with a unique mission title, transcription, and YouTube video link. 

Using React's `useEffect`, the page is then reloaded using this newly retrieved information, populating the page with the mission's title, the transcription, and the link to the original YouTube video the transcription was obtained from. 

The YouTube video is embedded using the original URL making sure that FREAKBAiT still gets any views that come about from this application, and also allows users to view the contents of the mission without audio if necessary.

Users can also use the "Regenerate Mission Parameters" button at the bottom of the page to change their session token, in case they'd like a different challenge than the one they're given. 

Since the session token is only regenerated if it does not already exist, or at the user's request, the same token is used from that point forward to make requests to the API. Additionally, because this token is randomly generated, no two users are likely to have the same one, and as a result, the missions they receive will be different as well, providing each user a unique experience.

A user can regenerate the token until their mission for the day matches another user's, but the next day they will be different again because the tokens being used are inherently different, aided in part by the large mission pool the result is being retrieved from.

Finally, near the mission title is a simple Dark vs. Light mode button, a preference that is also persisted within localStorage.

> The application is also a Progressive Web Application (PWA). I did this so the user could save the website to their home application, allowing for a somewhat native app experience, but without having to develop for different platforms or having to use a framework like React Native while keeping accessible from all devices.

## Deployment
Each component of the application also has a simple Dockerfile I created to make it very easy to containerize and deploy the application for deployment. 

The containers I made are then put on DockerHub, after which I pull them down to my VPS and deploy them with a simple Docker Compose file, though they would also work in a Kubernetes (K8s) environment.

## Further Improvements
All components of this app were made in a burst of inspration over the course of ~24 hours, and as a result it has a number of pitfalls.

The transcription is a script, it simply runs and then ends. This is not inherently wrong, but a different method would've allowed for a number of additions. It also currently just overwrites whatever is saved locally when rerun. Due to the nature of a Large Language Model like GPT-4, the resulting transcription can differ slightly between requests. If I were to do this again, I'd make it a constantly running service that checks for updates to the YouTube playlist, before transcribing the new additions, while also ignoring existing locally stored transcriptions. 

In a Kubernetes environment, I'd likely keep it as a script run by a CronJob deployment.

As for the communication between the frontend and the API, I'd have the frontend store whether or not the session token has changed since the last time the server was queried for the mission, and then cache the mission upon query. This prevents needing to make repeated queries for the same information to the server in the same day.