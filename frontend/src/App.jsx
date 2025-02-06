import React from 'react';
import Signup from './components/SignUp';
import CourseDetail from './components/CourseDetail';
import Login from './components/Login'
import CourseSearch from './components/CourseSearch';
// import AddComment from './components/AddComment'
import MyCoursesList from './components/MyCoursesList'
import AddCourse from './components/AddCourse'
import UpdateCourse from './components/UpdateCourse'
import Home from './components/Home';
import ContainerList from './components/ContainerList'
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";

const App = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path='/' element={<Navigate to="/login"/>} />
                <Route path='/home' element={<Home />}/>
                <Route path='/login' element={<Login />}/>
                <Route path='/signup' element={<Signup />} />
                <Route path="/courses/:id" element={<CourseDetail />} />
                <Route path="/courses/search" element={<CourseSearch />} />
                {/* <Route path="/comment" element={<AddComment />} /> */}
                <Route path="/mycourses" element={<MyCoursesList />} />
                <Route path="/create" element={<AddCourse/>} />
                <Route path="/containers" element={<ContainerList/>} />
                <Route path="/update/:id" element={<UpdateCourse/>} />
            </Routes> 
        </BrowserRouter>
    );
};

export default App;