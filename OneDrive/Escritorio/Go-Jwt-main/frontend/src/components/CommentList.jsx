import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import Comment from "./Comment.jsx";

function CommentsList() {
    const { courseId } = useParams();
    const [comments, setComments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch(`http://localhost:8080/comment/${courseId}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    credentials: 'include',
                });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data = await response.json();
                setComments(data || []); // Asegúrate de que `data` sea un array o un array vacío
                setError(null);
            } catch (error) {
                console.error("Error fetching comments:", error);
                setError('Error al cargar los datos. Por favor, intenta nuevamente.');
                setComments([]);
            } finally {
                setLoading(false);
            }
        }

        fetchData();
    }, [courseId]);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error fetching comments: {error}</div>;
    }

    return (
        <div className="comment-list">
            {comments && comments.length === 0 ? (
                <p>No hay comentarios para este curso.</p>
            ) : (
                comments.map(comment => (
                    <Comment key={comment.comment_id} comment={comment} />
                ))
            )}
        </div>
    );
}

export default CommentsList;
