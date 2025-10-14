import {Navigate, Outlet, useLocation, useNavigate} from "react-router";
import Header from "./header/Header.jsx";
import useAuth from "../hooks/useAuth.jsx";

const RequiredAUth = () => {
    const {auth} = useAuth();
    const location = useLocation();
    const navigate = useNavigate();
    return auth ? (
            <Outlet />
    ):  (
        Navigate({to: '/login', state: {from: location}, replace: true})
    )
}

export default RequiredAUth;