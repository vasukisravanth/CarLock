import axios from 'axios';

const apiClient = axios.create({
    baseURL: 'http://localhost:8080/api', // Adjust the base URL as needed
    timeout: 1000,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const lockCar = async (carId) => {
    const response = await apiClient.post(`/lock/${carId}`);
    return response.data;
};

export const unlockCar = async (carId) => {
    const response = await apiClient.post(`/unlock/${carId}`);
    return response.data;
};

export const getLockStatus = async (carId) => {
    const response = await apiClient.get(`/status/${carId}`);
    return response.data;
};