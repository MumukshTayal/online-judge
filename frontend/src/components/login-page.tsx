import { useState, useEffect } from 'react';
import { Button } from "@/components/ui/button";
import { jwtDecode } from "jwt-decode";
import { GoogleLogin } from '@react-oauth/google';

export function LoginPage() {
  const [jwtToken, setJwtToken] = useState('');

  // load the jwttoken from localstorage
  useEffect(() => {
    const storedToken = localStorage.getItem('jwtToken');
    if (storedToken) {
      setJwtToken(storedToken);
    }
  }, []);

  // verify creds and add in db if not added already
  const handleLoginSuccess = async (credentialResponse) => {
    const idToken = credentialResponse.credential;
    localStorage.setItem('jwtToken', idToken);
    setJwtToken(idToken);

    const response = await fetch('http://localhost:8080/login', {
      method: 'POST',
      header: {
        'Authorization': `Bearer ${jwtToken}`,
        `Content-Type`: 'application/json'
      }
    });  

    if (response.ok) {
      const data = await response.json();
      console.log('Login status:', data);
    } else {
      console.log('Login request failed');
    }
  };

  const handleLoginError = () => {
    console.log('Login Failed');
  };

  // logout 
  const handleLogout = () => {
    localStorage.removeItem('jwtToken');
    setJwtToken('');
  };

  return (
    <div>
      <GoogleLogin
        onSuccess={handleLoginSuccess}
        onError={handleLoginError}
        useOneTap
      />
      <div>
        <p>JWT Token: {jwtToken}</p>
        <Button onClick={handleLogout}>Logout</Button>
      </div>
    </div>
  );
}
