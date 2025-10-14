import useAxiosPrivate from "../../hooks/useAxiosPrivate.jsx";
import {useEffect, useState} from "react";
import {Spinner} from "react-bootstrap";
import Movies from "../movies/Movies.jsx";



const Recommended = () => {
    const [movies , setMovies] = useState([]);
    const [loading , setLoading] = useState(false)
    const [message , setMessage] = useState(null);
    const axiosPrivate = useAxiosPrivate();

    useEffect(() => {
        const getMovies = async () => {
            try {
                setLoading(true);
                const {data} = await axiosPrivate.get('/movies/recommended');
                console.log("DATA:", data)
                if(!data.movies || data.movies.length === 0){
                    setMessage('There are currently no recommended movies')
                }
                setMovies(data.movies);
            } catch (error) {
                console.error("Error fetching recommended movies", error)
            } finally {
                setLoading(false);
            }
        }
        getMovies();
    }, [])
    return (
        <>
            {loading ? (
                    <Spinner animation="border" role="status">
                        <span className="visually-hidden">Loading...</span>
                    </Spinner>
                ):
                <Movies movies={movies} message={message}/>
            }
        </>
    );
};

export default Recommended;