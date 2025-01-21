import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import './styles/CourseDetail.css';

const CourseDetail = () => {
    const [course, setCourse] = useState(null);

    const handleEnroll = async () => {
        try {
            // Obtén el courseID del localStorage
            const storedCourseID = localStorage.getItem('courseID');
            if (!storedCourseID) {
                console.error("courseID no encontrado en localStorage.");
                return;
            }
            const storedUserID = localStorage.getItem('userID');
            if (!storedUserID) {
                console.error("userID no encontrado en localStorage.");
                return;
            }

            const bodyContent = { user_id: storedUserID };
            console.log("Enviando body:", bodyContent);  // Imprimir el objeto que se enviará
            const response = await fetch(`http://localhost:8081/inscriptions/courses/${storedCourseID}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(bodyContent),
                credentials: 'include'
            });
            
            console.log(response)

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

    useEffect(() => {
        // Obtén el courseID del localStorage
        const storedCourseID = localStorage.getItem('courseID');

        if (!storedCourseID) {
            console.error("courseID no encontrado en localStorage.");
            return;
        }

        async function fetchCourse() {
            console.log("Fetching course with ID:", storedCourseID);
            try {
                const response = await fetch(`http://localhost:8081/courses/${storedCourseID}`, {
                    credentials: 'include',
                });
                const data = await response.json();
                console.log("Data received from API:", data);

                if (data && data._id) {
                    setCourse(data);
                } else {
                    console.error("error fetching the course");
                }
            } catch (error) {
                console.error('Error fetching course:', error);
            }
        }

        fetchCourse();
    }, []);

    if (!course) {
        return <div className="course-detail">Loading...</div>;
    }

    return (
        <div>
            <div className="course-detail">
                <h2>Detalles del Curso</h2>
                <p><strong>Nombre:</strong> {course.name}</p>
                <p><strong>Descripción:</strong> {course.description}</p>
                {/* <p><strong>Profesor:</strong> {course.professor}</p> */}
                <p><strong>Duración:</strong> {course.duration}</p>
                <p><strong>Requisito:</strong> {course.requirement}</p>
                <p> <img src={course.image_url} alt="No hay foto" /> </p>
                
                {/* Verifica si la disponibilidad es mayor a 0 */}
                {course.availability > 0 ? (
                    <button className="enroll-button" onClick={handleEnroll}>Inscribirse</button>
                ) : (
                    <p className="full-course-message">Curso completo</p>
                )}
                
                <Link to="/home">
                    <button className="back-button">Volver al Inicio</button>
                </Link>
            </div>
            <hr />
        </div>
    );
};

export default CourseDetail;
