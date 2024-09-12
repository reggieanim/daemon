import './style.css';
import './app.css';

import logo from './assets/images/test.png';
import { StartService, StopService } from '../wailsjs/go/main/App';

// Add SVG icon for the button
const playIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="feather feather-play" viewBox="0 0 24 24"><path d="M5 3L19 12 5 21 5 3z"/></svg>`;
const stopIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="feather feather-stop" viewBox="0 0 24 24"><rect width="18" height="18" x="3" y="3" rx="2" ry="2"/></svg>`;

document.querySelector('#app').innerHTML = `
    <img id="logo" class="logo">
    <div class="result" id="result">Service is currently stopped</div>
    <div class="input-box">
        <button id="toggleButton" class="btn">
            ${playIcon} Start
        </button>
    </div>
`;

// Set logo image
document.getElementById('logo').src = logo;

let resultElement = document.getElementById("result");
let buttonElement = document.getElementById("toggleButton");
let isServiceRunning = false; // Tracks if the service is running

const toggleService = () => {
    if (isServiceRunning) {
        StopService()
            .then(() => {
                isServiceRunning = false;
                resultElement.innerText = "Service is currently stopped";
                buttonElement.innerHTML = `${playIcon} Start Service`;
            })
            .catch((err) => {
                console.error("Error stopping service:", err);
            });
    } else {
        // Start the service
        StartService()
            .then(() => {
                isServiceRunning = true;
                resultElement.innerText = "Service is currently running";
                buttonElement.innerHTML = `${stopIcon} Stop Service`;
            })
            .catch((err) => {
                console.error("Error starting service:", err);
            });
    }
};

buttonElement.addEventListener("click", toggleService);
