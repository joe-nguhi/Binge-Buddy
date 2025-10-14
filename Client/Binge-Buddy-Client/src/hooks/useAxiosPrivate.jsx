import useAuth from "./useAuth.jsx";
import axios from "axios";
const apiUrl = import.meta.env.VITE_BASE_API_URL;


const useAxiosPrivate = () => {
    const {auth} = useAuth();
    const axiosAuth = axios.create({
        baseURL: apiUrl,
        headers: {'Content-Type': 'application/json'},
        withCredentials: true
    });

    axiosAuth.interceptors.request.use(
        config => {
            if (auth) {
                config.headers['Authorization'] = `Bearer ${auth.token}`;
            }
            return config;
        }
    )

    return axiosAuth;
}

export default useAxiosPrivate;