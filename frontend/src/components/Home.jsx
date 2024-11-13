import React, { useEffect, useState } from 'react';
import CourseList from './CourseList.jsx';
import './styles/Home.css';
import { Link, useNavigate  } from 'react-router-dom';

const Home = () => {
    const [email, setEmail] = useState('');
    const [isAdmin, setIsAdmin] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
        // Obtener email y tipo de usuario almacenados en localStorage al cargar la página
        const storedEmail = localStorage.getItem('email');
        const storedIsAdmin = localStorage.getItem('isAdmin') === 'admin'; // Convertir a booleano

        if (storedEmail) {
            setEmail(storedEmail);
            setIsAdmin(storedIsAdmin);
        } else {
            // Si no hay email almacenado, redirigir a la página de registro
            navigate('/login');
        }
    }, [navigate]);

    return (
        <div className="home-container">
            <header>
                <h1 id="title">I N I C I O</h1>
                <div id="infoUser" className={isAdmin ? 'admin' : 'alumno'}>
                    {email && <p>Usuario: {email}</p>}
                    {isAdmin ? <p>[ADMIN]</p> : <p>[ALUMNO]</p>}
                </div>
                <Link to='/mycourses'>
                    <button className="mycourses-button">Mis Cursos</button>
                </Link>
            </header>
            <hr />
            <div className="courses-section">
                <h2>Cursos Disponibles</h2>
                <Link to='/courses/search'>
                    <button className="search-button">Buscar Cursos</button>
                </Link> 
                {/* {isAdmin ? <Link to='/create'> <button className='createcourse-button'>Crear Curso</button> </Link> : <p></p>} */}
                
            </div>
            <hr />
            <CourseList />
        </div>
    );
};

export default Home;