import React, { useEffect, useState } from "react";
import { Link } from 'react-router-dom';
import MyCourse from "./MyCourse.jsx";
import './styles/MyCoursesList.css';

function MyCourses() {
    const [courses, setCourses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch(`http://localhost:8080/mycourses`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    credentials: 'include',
                });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data = await response.json();
                setCourses(data || []); // Asegúrate de que `data` sea un array o un array vacío
                setError(null);
            } catch (error) {
                console.error("Error fetching courses:", error);
                setError('Error al cargar los datos. Por favor, intenta nuevamente.');
                setCourses([]);
            } finally {
                setLoading(false);
            }
        }

        fetchData();
    }, []);

    if (loading) {
        return <div className="course-list">Cargando...</div>;
    }

    if (error) {
        return <div className="course-list">{error}</div>;
    }

    return (
        <div>
            <Link to="/home">
                <button className="back-button">Volver al Inicio</button>
            </Link>
            <header>
                <h1>Mis Cursos</h1>
            </header>
            <hr />
            <div className="course-list">
                {
                    courses.length > 0 ? (
                        courses.map(course => (
                            <MyCourse key={course.courseID} course={course} />
                        ))
                    ) : (
                        <div>No tienes cursos registrados.</div>
                    )
                }
            </div>
        </div>
    );
}

export default MyCourses;
