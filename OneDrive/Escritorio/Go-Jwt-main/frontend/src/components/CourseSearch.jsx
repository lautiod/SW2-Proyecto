import React, { useState } from 'react';
import Course from "./Course.jsx";
import { Link } from 'react-router-dom';
import './styles/CourseSearch.css';

const CourseSearch = () => {
    const [searchQuery, setSearchQuery] = useState('');
    const [coursesResult, setCoursesResult] = useState([]);
    const [error, setError] = useState('');
    const [searchOnce, setSearchOnce] = useState(false);

    const handleSearchChange = (event) => {
        setSearchQuery(event.target.value);
    };

    const handleSearchSubmit = async (event) => {
        event.preventDefault();

        try {
            const response = await fetch(`http://localhost:8080/courses/search?q=${searchQuery}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });
            const data = await response.json();

            setCoursesResult(data);
            setError('');
            setSearchOnce(true);
            console.log(data);
        } catch (error) {
            console.error("Error fetching courses:", error);
            setError('Error al cargar los datos. Por favor, intenta nuevamente.');
            setCoursesResult([]);
        }
    };

    return (
        <div>
            <h1>Busqueda de Cursos</h1>
            <Link to="/home">
                <button className="back-button">Volver al Inicio</button>
            </Link>
            <form onSubmit={handleSearchSubmit} className="search-form">
                <input
                    type="text"
                    placeholder="Buscar cursos..."
                    value={searchQuery}
                    onChange={handleSearchChange}
                    className="search-bar"
                />
                <button type="submit" className="search-button">Buscar</button>
            </form>
            {error && <p className="error-message">{error}</p>}
            {(!coursesResult || coursesResult.length === 0) && searchOnce && !error && (
                <p className="no-results">No se han encontrado resultados.</p>
            )}
            <div className="course-list">
                {coursesResult && coursesResult.map(course => (
                    <Course key={course.courseID} course={course} />
                ))}
            </div>
        </div>
    );
};

export default CourseSearch;
