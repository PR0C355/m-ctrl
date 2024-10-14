import React, { useState, useEffect } from 'react';
import Mission from './Mission';

const App = () => {

  // State to hold the mission
  const [mission, setMission] = useState({
    title: 'Initial Mission',
    instructions: 'These are the initial instructions.',
    id: 'Initial Video'
  });

  // Function to generate a new token
  function generateNewToken() {
    // Generate a UUID v4
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
      var r = Math.random() * 16 | 0;
      var v = c === 'x' ? r : (r & 0x3) | 0x8;
      return v.toString(16);
    });
  }

  // Function to regenerate the token
  function regenerateToken() {
    const newToken = generateNewToken();
    localStorage.setItem('sessionToken', newToken);

    window.location.reload(); 
  }

  // Function to generate a daily hash
  function generateDailyHash() {
    // Check for existing session token, if not found, generate a new one
    let sessionToken = localStorage.getItem('sessionToken');
    if (!sessionToken) {
      sessionToken = generateNewToken();
      localStorage.setItem('sessionToken', sessionToken);
    }

    const currentDate = new Date().toISOString().split('T')[0];
    const combinedString = sessionToken + '-' + currentDate;
    return combinedString
  }

  const generateNewMission = () => {
    fetch('http://localhost:11434/missions/get/random')
      .then((res) => {
        return res.json();
      })
      .then((data) => {
        console.log(data);
        const randomMission = {'title': data.title , 'instructions': data.text, 'id': data.id};
        setMission(randomMission);
      });
  };

  function fetchUniqueResult() {
    const dailyHash = generateDailyHash();
    
    // THis is the fetch for the local server
    // fetch('http://localhost:11434/missions/get/unique?token=' + dailyHash)

    // This is the fetch for the remote server
    fetch('https://mission.tumi.dev/missions/get/unique?token=' + dailyHash)
      .then((res) => {
        return res.json();
      })
      .then((data) => {
        // console.log(data);
        const uniqueMission = {'title': data.title , 'instructions': data.text, 'id': data.id};
        setMission(uniqueMission);
      });
  }


  function Main() {
    useEffect(() => {
      // Fetch unique user & day mission
      fetchUniqueResult();
    }, []); // Empty dependency array
  
    return (
      <Mission
      title={mission.title}
      instructions={mission.instructions}
      id={mission.id}
      onNewMission={regenerateToken}
    />
    );
  }
  
  return (Main());
};

export default App;