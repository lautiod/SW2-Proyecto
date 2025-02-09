import React from 'react';
import './styles/Container.css';

function Container({ container }) {
    return (
        <div className="course-card">
            <div className="course-content">
                <h2>
                    {container.name} <span className="container-id">id: {container.container_id}</span>
                </h2>
            </div>
        </div>
    );
}

export default Container;
