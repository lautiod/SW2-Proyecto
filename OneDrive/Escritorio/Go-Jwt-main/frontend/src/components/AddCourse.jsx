import React, { useState } from 'react';
import './styles/AddCourse.css'; 
import { Link } from 'react-router-dom';

const AddCourse = () => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [category, setCategory] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        const newCourse = {
            name: name,
            description: description,
            category: category
        };

        console.log(newCourse); 
        try {
            const response = await fetch('http://localhost:8080/createcourse', {
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
                setCategory(''); 
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
                    value={category}
                    onChange={(e) => setCategory(e.target.value)}
                    placeholder="Categoría del curso"
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
