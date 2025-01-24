import React, { useState } from 'react';
import { Navigate, Link } from 'react-router-dom';
import './styles/Login.css'; // Importar el archivo de estilos

const Login = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [redirect, setRedirect] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch('http://localhost:8080/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password }),
                credentials: 'include'
            });
            const data = await response.json();
            console.log(data);

            if (response.ok) {
                localStorage.setItem('email', email);
                localStorage.setItem('userID', data.id);
                localStorage.setItem('isAdmin', data.is_admin ? 'admin' : 'alumno');
                alert('Login exitoso');
                setRedirect(true);
            } else {
                alert('Error en el login: ' + data.message);
            }
        } catch (error) {
            console.error('Error al iniciar sesión:', error);
            alert('Error al iniciar sesión');
        }
    };

    if (redirect) {
        return <Navigate to="/home" />;
    }

    return (
        <div className="login-container">
            <h1 className="login-title">Login</h1>
            <form onSubmit={handleSubmit} className="login-form">
                <div className="form-group">
                    <label>Email:</label>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
                </div>
                <div className="form-group">
                    <label>Password:</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </div>
                <button type="submit" className="login-button">Login</button>
            </form>
            <Link to="/signup">
                <button className="signup-button">Crear Cuenta</button>
            </Link>
        </div>
    );
};

export default Login;