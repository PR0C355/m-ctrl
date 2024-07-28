import React, { useState } from 'react';
import YoutubeEmbed from "./YoutubeEmbed";
import { DarkModeSwitch } from 'react-toggle-dark-mode';
import "./App.css";


function MissionTitle({ title }) {
    return (
        <h1 className="mission-title">{title}</h1>
    );
}  

function MissionInstructions({ instructions }) {
    return (
        <>
        <h3 className="mission-instructions">{instructions}</h3>
        </>
        // <p className="mission-instructions">{instructions}</p>
    );
}

const Mission = ({ title, instructions, id, onNewMission }) => {
    const [isDarkMode, setIsDarkMode] = useState(localStorage.getItem('darkMode') === 'true' ? true : false);

    
    const lightTheme = {
        backgroundColor: '#f0f0f0',
        color: '#333',
    };
    
    const darkTheme = {
        backgroundColor: '#000000',
        color: '#ecf0f1',
    };

    const toggleDarkMode = () => {
        setIsDarkMode(!isDarkMode);
        localStorage.setItem('darkMode', !isDarkMode);
    };

    function Header() {
        return (
          <header>
            <div style={{display: 'flex', alignItems: 'center', justifyContent: 'center'}}>
                <h1 className="head">Mission Control</h1>
                <div style={{'marginLeft': '10px'}}>
                <DarkModeSwitch
                checked={isDarkMode}
                onChange={toggleDarkMode}
                size={40}
                />
                </div>
            </div>
          </header>
        );
    }
    
        
    return (
        
        <div className="app" style={isDarkMode ? darkTheme : lightTheme}>
            <div className="mission-container">
                <Header />
                <br />
                <hr />
                <br />
                <MissionTitle title={title} />
                <MissionInstructions instructions={instructions} />
                <YoutubeEmbed embedId={id} />
                <br />
                <button className={isDarkMode ? "new-mission-btn-dark" : "new-mission-btn-light"} onClick={onNewMission}>
                Regenerate Mission Parameters
                </button>
            </div>
        </div>
    );
};

export default Mission;