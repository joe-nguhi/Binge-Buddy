import axios from "axios";

const apiUrl = import.meta.env.VITE_BASE_API_URL;

export default axios.create({
    baseURL: apiUrl,
    headers: {
        "Content-type": "application/json"
    }
})