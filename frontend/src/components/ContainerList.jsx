import React, { useEffect, useState } from "react";
import { Link } from 'react-router-dom';
import Container from "./Container.jsx";
import './styles/ContainerList.css';

function ContainerList() {
    const [containers, setContainers] = useState([]);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch(`http://localhost:8080/users/admin`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include'
                })
                const data = await response.json();
                
                setContainers(data);
                console.log(data);
            } catch (error) {
                console.error("Error fetching containers:", error);
            }
        }

        fetchData();
    }, []);

    return (
        <>
          <h2 className="admin-title">Admin</h2>
          <hr />
          <div className="course-list">
            {containers.map(container => (
              <Container key={container.container_id} container={container} />
            ))}
          </div>
          <Link to="/home">
                          <button className="back-button">Volver al Inicio</button>
                      </Link>
        </>
      );      
}

export default ContainerList;