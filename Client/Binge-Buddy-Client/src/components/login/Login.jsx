import Client from "../../api/axios-config.js"
import {useState} from "react";
import {Link, useLocation, useNavigate} from "react-router";
import {Button, Container, Form} from "react-bootstrap";

const Login = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            const {data} = await Client.post('/login', {email, password});
            localStorage.setItem('token', data.user.token);
            localStorage.setItem('user', JSON.stringify(data.user));
            navigate(location.state?.from || '/', {replace: true});
        } catch (error) {
            setError(error.response.data.message);
        } finally {
            setLoading(false);
        }
    };
    return (
        <Container className="login-cotainer d-flex align-items-center justify-content-center min-vh-100">
            <div className="login-card shadow p-4 rounded bg-white" style={{maxWidth: 400, width: '100%'}}>
                <div className="text-center mb-4">
                    <h2 className="fw-bold">Sign In</h2>
                    <p className="text-muted">Welcome back! Please login to your account</p>
                </div>
                {error && <div className="alert alert-danger py-2">{error}</div>}
                <Form onSubmit={handleSubmit}>
                    <Form.Group controlId="email" className="mb-3">
                        <Form.Label>Email address</Form.Label>
                        <Form.Control
                            type="email"
                            placeholder="Enter email"
                            autoComplete="username"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            autoFocus
                        />
                    </Form.Group>

                    <Form.Group controlId="password" className="mb-3">
                        <Form.Label>Password</Form.Label>
                        <Form.Control
                            type="password"
                            placeholder="Password"
                            autoComplete="current-password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                    </Form.Group>


                    <Button variant="primary" type="submit" className="w-100 mb-2" disabled={loading}
                            style={{fontWeight: 600, letterSpacing: 1}}>
                        {loading ? (
                            <>
                            <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true">
                                Loggin in...
                            </span>
                            </>
                        ) : 'Login'}
                    </Button>
                </Form>
                <div className="text-center mt-3">
                    <p className="text-muted">Don't have an account? <Link to="/register">Sign up</Link></p>
                </div>
            </div>

        </Container>
    )
}

export default Login;