import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import './styles/Course.css';

function Course({ course }) {
    const [isAdmin, setIsAdmin] = useState(false); // Estado para isAdmin

    useEffect(() => {
        // Obtén isAdmin desde localStorage
        const storedIsAdmin = localStorage.getItem('isAdmin') === 'admin'; // Convertir a booleano
        setIsAdmin(storedIsAdmin);
    }, []);

    const handleButton = () => {
        // Guarda el course._id en el localStorage
        localStorage.setItem('courseID', course._id);
        // Redirige a la página de detalles del curso
        // navigate(`/courses/${course._id}`);
    };

    return (
        <div className="course-card">
            <div className="course-content">
                <h2>
                    {course.name}
                    <img src={course.image_url} alt="No hay imagen del curso" />
                </h2>
                <Link to={`/courses/${course._id}`}>
                    <button className="info-button" onClick={handleButton}>Información</button>
                </Link>
                {isAdmin ? <Link to={`/update/${course._id}`}> <button className='updatecourse-button' onClick={handleButton}>Actualizar Curso</button> </Link> : <p></p>}

            </div>
        </div>
    );
}

export default Course;