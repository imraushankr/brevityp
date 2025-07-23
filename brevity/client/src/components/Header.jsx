import React from 'react';
import Logo from './Logo';
import "../styles/Headers.css";

const Header = () => {
  return (
    <header className='h_section'>
      <Logo />
      <nav className='nav_section'>
        <a href="#">Home</a>
        <a href="#">About</a>
        <a href="#">Services</a>
        <a href="#">Contact</a>
      </nav>
    </header>
  );
};

export default Header;