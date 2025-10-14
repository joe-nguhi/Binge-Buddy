import {Button, Container, Nav, Navbar, NavDropdown} from "react-bootstrap";
import {NavLink, useNavigate} from "react-router";
import {useState} from "react";

const Header = () => {
    const navigate = useNavigate();
    const [auth, setAuth] = useState(false);
    return (
        <>
            <Navbar bg="dark" variant="dark" expand="lg" fixed="top" className="shadow-sm">
                <Container>
                    <Navbar.Brand href="#home">Binge Buddy</Navbar.Brand>
                    <Navbar.Toggle aria-controls="basic-navbar-nav"/>
                    <Navbar.Collapse id="basic-navbar-nav">
                        <Nav className="me-auto">
                            <Nav.Link as={NavLink} to="/">Home</Nav.Link>
                            <Nav.Link as={NavLink} to="/recommended">Recommended</Nav.Link>
                        </Nav>
                        <Nav className="ms-auto align-items-center">
                            {auth ? (
                                <>
                            <span>
                                Hello, {auth.username}
                            </span>
                                    <Button variant="outline-danger" size="sm" onClick={() => {
                                        localStorage.removeItem("token");
                                        setAuth(false);
                                        navigate("/login");
                                    }}>Logout
                                    </Button>
                                </>
                            ) : (
                                <>
                                    <Button className="me-2" variant="outline-info" size="sm"
                                            onClick={() => navigate("/login")}>
                                        Login
                                    </Button>
                                    <Button className="me-2" variant="outline-info" size="sm"
                                            onClick={() => navigate("/register")}>
                                        Register
                                    </Button>
                                </>
                            )}
                        </Nav>
                    </Navbar.Collapse>
                </Container>
            </Navbar>
        </>
    )
}

export default Header;