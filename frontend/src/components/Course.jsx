import React from 'react';
import { Link } from 'react-router-dom';
import './styles/Course.css';

function Course({ course }) {
    // const handleEnroll = async () => {
    //     try {
    //         const response = await fetch('http://localhost:8080/subscribe', {
    //             method: 'POST',
    //             headers: {
    //                 'Content-Type': 'application/json',
    //             },
    //             body: JSON.stringify({ courseID: course.courseID}), // Adjust this based on how you identify the user
    //             credentials: 'include'
    //         });

    //         if (response.ok) {
    //             alert('Inscripción exitosa!');
    //         } else {
    //             alert('Error al inscribirse en el curso.');
    //         }
    //     } catch (error) {
    //         console.error('Error:', error);
    //         alert('Error al inscribirse en el curso.');
    //     }
    // };

    const handleMoreInfoClick = () => {
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
                    <button className="info-button" onClick={handleMoreInfoClick}>Información</button>
                </Link>
                {/*<button className="enroll-button" onClick={handleEnroll}>Inscribirse</button> */}
            </div>
        </div>
    );
}

export default Course;