import React from 'react';
import './styles/Comment.css';

function Comment({ comment }) {
    return (
        <div className="comment-card">
            <div className="comment-content">
                <h3>Usuario: {comment.email}</h3>
                <p>Comentario: {comment.content}</p>
            </div>
        </div>
    );
}

export default Comment;
