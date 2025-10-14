import Client from '../../api/axios-config.js'
import {useEffect, useState} from "react";
import {Button, Container, Form} from "react-bootstrap";
import {useNavigate} from "react-router";

const Register = () => {
    const navigate = useNavigate();
    const [firstName, setFirstName] = useState('');
    const [lastName, setLastName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [favoriteGenres, setFavoriteGenres] = useState([]);
    const [genres, setGenres] = useState([]);

    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        const getGenres = async () => {
            setLoading(true)
            try {
                const {data} = await Client.get('/movies/genres');
                setGenres(data.genres);
            } catch (error) {
                console.log("Error fecthing movie genres:", error.message)
            }
            setLoading(false)
        };
        getGenres();
    }, []);

    const handleGenreChange = (e) => {
        const selectedGenres = Array.from(e.target.selectedOptions, option => ({
            genre_id: parseInt(option.value),
            genre_name: option.label
        }));
        setFavoriteGenres(selectedGenres);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);
        const defaultRole = 'USER';

        if (password !== confirmPassword) {
            setError('Passwords do not match');
            return;
        }
        setLoading(true)

        try {
            const payload = {
                first_name: firstName,
                last_name: lastName,
                email,
                password,
                role: defaultRole,
                favorite_genres: favoriteGenres
            };

            const {data} = await Client.post('/register', payload);
            if (data.error){
                setError(data.error);
            }
            navigate("/login", {replace:true})
        }catch (_) {
            setError("Registration Failed. Please try again later")
        }finally {
            setLoading(false)
        }
    };


    return (
        <Container className="login-cotainer d-flex align-items-center justify-content-center min-vh-100">
            <div className="login-card shadow p-4 rounded bg-white" style={{maxWidth:400, width:'100%'}}>
                <div className="text-center mb-4">
                    <h2 className="fw-bold">Register</h2>
                    <p className="text-muted">Create an account to start streaming  your favorite movies</p>
                    {error && <div className="alert alert-danger py-2">{error}</div>}
                </div>
                <Form onSubmit={handleSubmit}>
                    <Form.Group>
                        <Form.Label>First Name</Form.Label>
                        <Form.Control
                            type="text"
                            placeholder="Enter first name"
                            value={firstName}
                            onChange={(e) => setFirstName(e.target.value)}
                            required
                        />
                    </Form.Group>
                    <Form.Group>
                        <Form.Label>Last Name</Form.Label>
                        <Form.Control
                            type="text"
                            placeholder="Enter last name"
                            value={lastName}
                            onChange={(e) => setLastName(e.target.value)}
                            required
                        />
                    </Form.Group>

                    <Form.Group className="mb-3">
                        <Form.Label>Email address</Form.Label>
                        <Form.Control
                            type="email"
                            placeholder="Enter email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                        />
                        <Form.Text className="text-muted">
                            We'll never share your email with anyone else.
                        </Form.Text>
                    </Form.Group>

                    <Form.Group className="mb-3" >
                        <Form.Label>Password</Form.Label>
                        <Form.Control
                            type="password"
                            placeholder="Password"
                            autoComplete="new-password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                    </Form.Group>


                    <Form.Group className="mb-3" >
                        <Form.Label>Confirm Password</Form.Label>
                        <Form.Control
                            type="password"
                            placeholder="Confirm Password"
                            autoComplete="new-password"
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                            required
                            isInvalid={!confirmPassword && confirmPassword !== password}
                        />
                        <Form.Control.Feedback type="invalid">
                            Passwords do not match
                        </Form.Control.Feedback>
                    </Form.Group>

                    <Form.Group>
                        <Form.Label>Favorite Genres</Form.Label>
                        <Form.Select
                            multiple
                            value={favoriteGenres.map((genre) => String(genre.genre_id))}
                            onChange={handleGenreChange}
                            required
                        >
                            {/*<option value="">Select a genre</option>*/}
                            {genres.map((genre) => (
                                <option key={genre.genre_id} value={genre.genre_id} label={genre.genre_name}>
                                    {genre.genre_name}
                                </option>
                            ))}
                        </Form.Select>
                        <Form.Text className="text-muted">
                            Hold Ctrl or Cmd to select multiple genres
                        </Form.Text>
                    </Form.Group>

                    <Button variant="primary" type="submit" className="w-100 mb-2" disabled={loading} style={{fontWeight:600, letterSpacing:1}} >
                        {loading ? (
                            <>
                            <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true">
                                Registering...
                            </span>
                            </>
                        ):'Register'}
                    </Button>
                </Form>
            </div>

        </Container>
    )
}

export default Register;