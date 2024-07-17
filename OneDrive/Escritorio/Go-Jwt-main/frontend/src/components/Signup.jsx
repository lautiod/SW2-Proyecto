import React, { useState } from 'react';
import { Navigate } from 'react-router-dom';
import './styles/Signup.css'; 

const Signup = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [isAdmin, setIsAdmin] = useState(false);
    const [redirect, setRedirect] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        const newUser = {
            email: email,
            password: password,
            isAdmin: isAdmin
        };
        try {
            const response = await fetch('http://localhost:8080/signup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(newUser)
            });

            if (response.ok) {
                alert("Registro exitoso");
                // Después del registro exitoso, establecer el estado para redirigir
                setRedirect(true);
            } else {
                const data = await response.json();
                alert("Error al registrarse: " + data.message);
            }
        } catch (error) {
            console.error('Error al registrar:', error);
            alert("Error al registrarse");
        }
    };

    if (redirect) {
        return <Navigate to="/login" />;
    }

    return (
        <div className="signup-container">
            <h1 className="signup-title">Signup</h1>
            <form onSubmit={handleSubmit} className="signup-form">
                <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="Email"
                    required
                />
                <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Password"
                    required
                />
                <div className="radio-group">
                    <input
                        type="radio"
                        id="alumno"
                        checked={!isAdmin}
                        onChange={() => setIsAdmin(false)}
                        required
                    />
                    <label htmlFor="alumno">Alumno</label>
                    <input
                        type="radio"
                        id="admin"
                        checked={isAdmin}
                        onChange={() => setIsAdmin(true)}
                        required
                    />
                    <label htmlFor="admin">Administrador</label>
                </div>
                <button type="submit" className="signup-button">Signup</button>
            </form>
        </div>
    );
};

export default Signup;
