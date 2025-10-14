import Client from '../../api/axios-config.js'
import {useEffect, useState} from "react";
import Movies from "../movies/Movies.jsx";
import {Spinner} from "react-bootstrap";

const Home = () => {
    const [movies, setMovies] = useState([]);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState("");

    useEffect(() => {
        const getMovies = async () => {
            setLoading(true);
            try {
                const {data} = await Client.get('/movies');

                if(data.movies.length === 0){
                    setMessage('There are currently no available movies')
                }
                setMovies(data.movies);
            } catch (error) {
                setMessage(error.message);
            }
            setLoading(false);
        };
        getMovies()
    }, []);

  return (
    <>
        {loading? (
                <Spinner animation="border" role="status">
                    <span className="visually-hidden">Loading...</span>
                </Spinner>

        ):
            <Movies movies={movies} message={message}/>
        }
    </>
  );
};

export default Home;