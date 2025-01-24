import React, { useState } from 'react';
import './styles/AddCourse.css'; 
import { Link } from 'react-router-dom';

const AddCourse = () => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [professor, setProfessor] = useState('');
    const [image_url, setImageURL] = useState('');
    const [duration, setDuration] = useState('');
    const [requirement, setRequirement] = useState('');
    const [availability, setAvailability] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();

         // Validar que duration sea un número válido
         if (isNaN(parseInt(duration))) {
            alert("La duración debe ser un número válido");
            return;
        }

        if (isNaN(parseInt(availability))) {
            alert("La disponibilidad debe ser un número válido");
            return;
        }

        const newCourse = {
            name: name,
            description: description,
            professor: professor,
            image_url: image_url,
            duration: parseInt(duration), // Convertir a número
            requirement: requirement,
            availability: parseInt(availability)
        };

        console.log(newCourse); 
        try {
            const response = await fetch('http://localhost:8081/courses', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(newCourse),
                credentials: 'include'
            });

            if (response.ok) {
                alert("Curso agregado correctamente");
                setName(''); 
                setDescription(''); 
                setProfessor(''); 
                setImageURL('');
                setDuration('');
                setRequirement('');
                setAvailability('');
            } else {
                const data = await response.json();
                alert("Error al agregar curso: " + data.message);
            }
        } catch (error) {
            console.error('Error al agregar curso:', error);
            alert("Error al agregar curso");
        }
    };

    return (
        <div className="course-container">
            <h1 className="course-title">Agregar Curso</h1>
            <form onSubmit={handleSubmit} className="course-form">
                <input
                    type="text" 
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="Nombre del curso"
                    required
                />
                <input
                    type="text" 
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    placeholder="Descripción del curso"
                    required
                />
                <input
                    type="text" 
                    value={professor}
                    onChange={(e) => setProfessor(e.target.value)}
                    placeholder="Profesor del curso"
                    required
                />
                <input
                    type="text" 
                    value={image_url}
                    onChange={(e) => setImageURL(e.target.value)}
                    placeholder="Imagen URL del curso"
                    required
                />
                <input
                    type="number" 
                    value={duration}
                    onChange={(e) => setDuration(e.target.value)}
                    placeholder="Duracion del curso"
                    required
                />
                <input
                    type="text" 
                    value={requirement}
                    onChange={(e) => setRequirement(e.target.value)}
                    placeholder="Requerimentos del curso"
                    required
                />
                <input
                    type="number" 
                    value={availability}
                    onChange={(e) => setAvailability(e.target.value)}
                    placeholder="Disponibilidad del curso"
                    required
                />
                <button type="submit" className="course-button">Agregar Curso</button>
            </form>
            <Link to="/home">
                <button className="back-button">Volver al Inicio</button>
            </Link>
        </div>
    );
};

export default AddCourse;