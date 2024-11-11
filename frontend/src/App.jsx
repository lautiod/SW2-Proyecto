import React from 'react';
// import Signup from './components/Signup';
// import CourseDetail from './components/CourseDetail';
// import Login from './components/Login';
// import CourseSearch from './components/CourseSearch';
// import AddComment from './components/AddComment'
// import MyCoursesList from './components/MyCoursesList'
// import AddCourse from './components/AddCourse'
import Home from './components/Home';
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";


const App = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path='/' element={<Navigate to="/home"/>} />
                <Route path='/home' element={<Home />}/>
                {/* <Route path='/login' element={<Login />}/>
                <Route path='/signup' element={<Signup />} />
                <Route path="/course/:courseId" element={<CourseDetail />} />
                <Route path="/course/search" element={<CourseSearch />} />
                <Route path="/comment" element={<AddComment />} />
                <Route path="/mycourses" element={<MyCoursesList />} />
                <Route path="/create" element={<AddCourse/>} /> */}
            </Routes> 
        </BrowserRouter>
    );
};

export default App;