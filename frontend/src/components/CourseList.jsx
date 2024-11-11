import React, { useEffect, useState } from "react";
import Course from "./Course.jsx";
import './styles/CourseList.css';

function CourseList() {
    const [courses, setCourses] = useState([]);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch('http://localhost:8081/courses');
                const data = await response.json();
                
                setCourses(data);
                console.log(data);
            } catch (error) {
                console.error("Error fetching courses:", error);
            }
        }

        fetchData();
    }, []);

    return (
        <div className="course-list">
            {
                courses.map(course => (
                    <Course key={course._id} course={course} />
                ))
            }
        </div>
    );
}

export default CourseList;