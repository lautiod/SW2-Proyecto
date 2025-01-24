import React, { useState } from 'react';
import './styles/AddCourse.css';
import { Link } from 'react-router-dom';

const UpdateCourse = () => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [professor, setProfessor] = useState('');
    const [image_url, setImageURL] = useState('');
    const [duration, setDuration] = useState('');
    const [requirement, setRequirement] = useState('');
    const [availability, setAvailability] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();

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

        // Crear objeto solo con los campos llenados
        const updatedCourse = {};
        if (name) updatedCourse.name = name;
        if (description) updatedCourse.description = description;
        if (professor) updatedCourse.professor = professor;
        if (image_url) updatedCourse.image_url = image_url;
        if (duration) {
            const dur = parseInt(duration);
            if (isNaN(dur) || dur < 0) {
                alert("La duracion debe ser un número válido mayor o igual a cero");
                return;
            }
            updatedCourse.duration = dur;
        }
        if (requirement) updatedCourse.requirement = requirement;
        if (availability) {
            const avail = parseInt(availability);
            if (isNaN(avail) || avail < 0) {
                alert("La disponibilidad debe ser un número válido mayor o igual a cero");
                return;
            }
            updatedCourse.availability = avail;
        }

        console.log(updatedCourse);
        try {
            const response = await fetch(`http://localhost:8081/courses/${storedCourseID}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(updatedCourse),
                credentials: 'include'
            });

            if (response.ok) {
                alert("Curso actualizado correctamente");
                setName('');
                setDescription('');
                setProfessor('');
                setImageURL('');
                setDuration('');
                setRequirement('');
                setAvailability('');
            } else {
                const data = await response.json();
                alert("Error al actualizar curso: " + data.message);
            }
        } catch (error) {
            console.error('Error al actualizar curso:', error);
            alert("Error al actualizar curso");
        }
    };

    return (
        <div className="course-container">
            <h1 className="course-title">Actualizar Curso</h1>
            <form onSubmit={handleSubmit} className="course-form">
                <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="Nombre del curso"
                />
                <input
                    type="text"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    placeholder="Descripción del curso"
                />
                <input
                    type="text"
                    value={professor}
                    onChange={(e) => setProfessor(e.target.value)}
                    placeholder="Profesor del curso"
                />
                <input
                    type="text"
                    value={image_url}
                    onChange={(e) => setImageURL(e.target.value)}
                    placeholder="Imagen URL del curso"
                />
                <input
                    type="number"
                    value={duration}
                    onChange={(e) => setDuration(e.target.value)}
                    placeholder="Duracion del curso"
                />
                <input
                    type="text"
                    value={requirement}
                    onChange={(e) => setRequirement(e.target.value)}
                    placeholder="Requerimentos del curso"
                />
                <input
                    type="number"
                    value={availability}
                    onChange={(e) => setAvailability(e.target.value)}
                    placeholder="Disponibilidad del curso"
                />
                <button type="submit" className="course-button">Actualizar Curso</button>
            </form>
            <Link to="/home">
                <button className="back-button">Volver al Inicio</button>
            </Link>
        </div>
    );
};

export default UpdateCourse;