import React, { useState } from 'react';
import axios from 'axios';

const LockControl: React.FC = () => {
    const [isLocked, setIsLocked] = useState<boolean>(false);

    const toggleLock = async () => {
        try {
            const response = await axios.post('/api/lock/toggle', { locked: !isLocked });
            setIsLocked(response.data.locked);
        } catch (error) {
            console.error('Error toggling lock:', error);
        }
    };

    return (
        <div>
            <h1>Car Lock Control</h1>
            <p>The car is currently {isLocked ? 'Locked' : 'Unlocked'}</p>
            <button onClick={toggleLock}>
                {isLocked ? 'Unlock Car' : 'Lock Car'}
            </button>
        </div>
    );
};

export default LockControl;