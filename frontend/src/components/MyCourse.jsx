import React from 'react';
import { Link } from 'react-router-dom';
import './styles/Course.css';

function MyCourse({ course }) {

    return (
        <div className="course-card">
            <div className="course-content">
                <h2>{course.name}</h2>
                <img src={course.image_url} alt="No hay imagen del curso" />
            </div>
        </div>
    );
}

export default MyCourse;