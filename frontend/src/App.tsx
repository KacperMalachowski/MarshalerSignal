import './App.css';
import {GetStations} from '../wailsjs/go/main/App'
import { useState } from 'react';

function App() {

    const [stations, setStations] = useState({})

    function getStations() {
        GetStations().then(result => setStations(result))
    }

    getStations()
    return (
        <div id="App">
            <pre>
                {JSON.stringify(stations)}
            </pre>
        </div>
    )
}

export default App
