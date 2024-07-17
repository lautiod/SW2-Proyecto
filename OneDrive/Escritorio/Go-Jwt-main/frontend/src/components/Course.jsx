import React from 'react';
import { Link } from 'react-router-dom';
import './styles/Course.css';

function Course({ course }) {
    const handleEnroll = async () => {
        try {
            const response = await fetch('http://localhost:8080/subscribe', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ courseID: course.courseID}), // Adjust this based on how you identify the user
                credentials: 'include'
            });

            if (response.ok) {
                alert('Inscripción exitosa!');
            } else {
                alert('Error al inscribirse en el curso.');
            }
        } catch (error) {
            console.error('Error:', error);
            alert('Error al inscribirse en el curso.');
        }
    };

    return (
        <div className="course-card">
            <div className="course-content">
                <h2>{course.name}</h2>
                <p>Categoria: {course.category}</p>
                <Link to={`/course/${course.courseID}`}>
                    <button className="info-button">Información</button>
                </Link>
                <button className="enroll-button" onClick={handleEnroll}>Inscribirse</button>
            </div>
        </div>
    );
}

export default Course;
