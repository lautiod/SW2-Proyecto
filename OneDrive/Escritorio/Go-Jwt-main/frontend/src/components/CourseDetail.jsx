import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import './styles/CourseDetail.css';

const CourseDetail = () => {
    const { courseId } = useParams();
    const [course, setCourse] = useState(null);

    const uploadFile = async (event) => {
        const file = event.target.files[0];

        if (file) {
            console.log('Nombre del archivo:', file.name);

            const reader = new FileReader();

            reader.onload = async () => {
                const result = reader.result.replace("data:", "").replace(/^.+,/, "");

                const response = await fetch('http://localhost:8080/addfile', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        file_name: file.name,
                        data: result,
                        course_id: course.courseID
                    }),
                    credentials: 'include'
                });

                if (response.ok) {
                    alert('Archivo subido con éxito!');
                } else {
                    alert('Error al subir el archivo.');
                }
            };

            reader.readAsDataURL(file);
        }
    };

    useEffect(() => {
        async function fetchCourse() {
            try {
                const response = await fetch(`http://localhost:8080/courses/${courseId}`, {
                    credentials: 'include',
                });
                const data = await response.json();
                setCourse(data);
                localStorage.setItem('courseID', data.courseID);
            } catch (error) {
                console.error('Error fetching course:', error);
            }
        }

        fetchCourse();
    }, [courseId]);

    if (!course) {
        return <div className="course-detail">Loading...</div>;
    }

    return (
        <div className="course-detail">
            <h2>Detalles del Curso</h2>
            <p><strong>Nombre:</strong> {course.name}</p>
            <p><strong>Descripción:</strong> {course.description}</p>
            <p><strong>Categoría:</strong> {course.category}</p>
            <div>
                <input type="file" onChange={uploadFile} />
            </div>
            <Link to="/home">
                <button className="back-button">Volver al Inicio</button>
            </Link>
        </div>
    );
};

export default CourseDetail;
