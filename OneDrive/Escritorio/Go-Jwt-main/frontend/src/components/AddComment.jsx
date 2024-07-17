import React, { useState } from 'react';
import './styles/AddComment.css'; 
import { Link } from 'react-router-dom';

const AddComment = () => {
    const [content, setContent] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        const storedCourseId = localStorage.getItem('courseID');
        const newComment = {
            course_id: Number(storedCourseId), 
            content: content
        };

        console.log(newComment); 
        try {
            const response = await fetch('http://localhost:8080/comment', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(newComment),
                credentials: 'include'
            });

            if (response.ok) {
                alert("Comentario agregado correctamente");
                setContent(''); // Limpiar el campo de comentario después de agregarlo
            } else {
                const data = await response.json();
                alert("Error al agregar comentario: " + data.message);
            }
        } catch (error) {
            console.error('Error al agregar comentario:', error);
            alert("Error al agregar comentario");
        }
    };

    return (
        <div className="comment-container">
            <h1 className="comment-title">Agregar Comentario</h1>
            <form onSubmit={handleSubmit} className="comment-form">
                <input
                    type="text" 
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="Escribe tu comentario"
                    required
                />
                <button type="submit" className="comment-button">Agregar Comentario</button>
            </form>
            <Link to={`/course/${localStorage.getItem('courseID')}`}>
                <button className="back-button">Volver al Detalle del Curso</button>
            </Link>
        </div>
    );
};

export default AddComment;
