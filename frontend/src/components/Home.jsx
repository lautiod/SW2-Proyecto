import React from 'react';
import CourseList from './CourseList.jsx';
import './styles/Home.css';

const Home = () => {
    return (
        <div className="home-container">
            <header>
                <h1 id="title">I N I C I O</h1>
                {/* <div id="infoUser" className={isAdmin ? 'admin' : 'alumno'}>
                    {email && <p>Usuario: {email}</p>}
                    {isAdmin ? <p>[ADMIN]</p> : <p>[ALUMNO]</p>}
                </div>
                <Link to='/mycourses'>
                    <button className="mycourses-button">Mis Cursos</button>
                </Link> */}
            </header>
            <hr />
            <div className="courses-section">
                <h2>Cursos Disponibles</h2>
                {/* <Link to='/course/search'>
                    <button className="search-button">Buscar Cursos</button>
                </Link> */}
                {/* {isAdmin ? <Link to='/create'> <button className='createcourse-button'>Crear Curso</button> </Link> : <p></p>} */}
                
            </div>
            <hr />
            <CourseList />
        </div>
    );
};

export default Home;