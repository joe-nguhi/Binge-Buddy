
import './App.css'
import Home from "./components/home/Home.jsx";
import Header from "./components/header/Header.jsx";
import {Route, Routes} from "react-router";
import Login from "./components/login/Login.jsx";
import Register from "./components/register/Register.jsx";
import Layout from "./components/Layout.jsx";
import RequiredAuth from "./components/RequiredAuth.jsx";
import Recommended from "./components/recommended/Recommended.jsx";

function App() {
  return (
      <>
          <Header />
          <Routes path="/" element={<Layout/>}>
              <Route index element={<Home/>}/>
              <Route path="/login" element={<Login/>}/>
              <Route path="/register" element={<Register/>}/>
              {/*Requires Auth*/}
              <Route element={<RequiredAuth/>}>
                  <Route path="/movies/recommended" element={<Recommended/>}/>
              </Route>

          </Routes>
      </>
  )
}

export default App;
