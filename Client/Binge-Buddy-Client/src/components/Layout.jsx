import {Outlet} from "react-router";
import Header from "./header/Header.jsx";

const Layout = () => {
    return (
        <main>
            <Outlet />
        </main>
    )
}

export default Layout;