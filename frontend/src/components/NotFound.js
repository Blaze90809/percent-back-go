import React from 'react';
import { Link } from 'react-router-dom';
import './NotFound.css';

const NotFound = () => {
  // Gentle background snow particles
  const snowflakes = [...Array(30)].map((_, i) => {
    const size = Math.random() * 4 + 2; // 2px to 6px
    const left = Math.random() * 100; // 0% to 100%
    const duration = Math.random() * 6 + 5; // 5s to 11s (gentle fall)
    const delay = Math.random() * 10; // 0s to 10s
    const opacity = Math.random() * 0.6 + 0.2; // 0.2 to 0.8
    
    return (
      <div
        key={i}
        className="snowflake"
        style={{
          left: `${left}%`,
          width: `${size}px`,
          height: `${size}px`,
          animationDuration: `${duration}s`,
          animationDelay: `${delay}s`,
          opacity: opacity,
        }}
      />
    );
  });

  return (
    <div className="not-found-container">
      {/* Dynamic Full-Canvas Vector Winter Poster */}
      <div className="poster-canvas">
        <svg viewBox="0 0 800 500" className="poster-svg" xmlns="http://www.w3.org/2000/svg">
          <defs>
            {/* Sky morning glow gradient */}
            <linearGradient id="sky-grad" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" stopColor="#e0f2f1" />
              <stop offset="60%" stopColor="#f5fbfb" />
              <stop offset="100%" stopColor="#ffffff" />
            </linearGradient>

            {/* Sun/Moon soft radial gradient */}
            <radialGradient id="sun-glow" cx="50%" cy="50%" r="50%">
              <stop offset="0%" stopColor="#fff9e6" stopOpacity="1" />
              <stop offset="40%" stopColor="#fffaec" stopOpacity="0.8" />
              <stop offset="100%" stopColor="#ffffff" stopOpacity="0" />
            </radialGradient>

            {/* Mountain shadows gradients */}
            <linearGradient id="mount-grad-far" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="#b2dfdb" />
              <stop offset="100%" stopColor="#80cbc4" />
            </linearGradient>
            
            <linearGradient id="mount-grad-mid" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="#80cbc4" />
              <stop offset="100%" stopColor="#4db6ac" />
            </linearGradient>

            <linearGradient id="mount-grad-close" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="#4db6ac" />
              <stop offset="100%" stopColor="#26a69a" />
            </linearGradient>

            {/* Foreground Skier Goggle Gradient */}
            <linearGradient id="goggle-fore" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stopColor="#00f2fe" />
              <stop offset="50%" stopColor="#4facfe" />
              <stop offset="100%" stopColor="#00ebfe" />
            </linearGradient>

            {/* Background Skier Goggle Gradient */}
            <linearGradient id="goggle-back" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stopColor="#ff9f43" />
              <stop offset="100%" stopColor="#ff3f34" />
            </linearGradient>

            {/* Lens Glare Clipping Paths */}
            <clipPath id="goggle-clip-fore">
              <path d="M 334 148 Q 346 142 358 148 Q 365 155 358 162 Q 346 166 334 162 Q 327 155 334 148 Z" />
            </clipPath>
            <clipPath id="goggle-clip-back">
              <path d="M 526 128 Q 534 123 542 128 Q 547 133 542 138 Q 534 141 526 138 Q 521 133 526 128 Z" />
            </clipPath>
          </defs>

          {/* 1. Sky */}
          <rect width="800" height="500" fill="url(#sky-grad)" />

          {/* 2. Morning Sun Glow */}
          <circle cx="620" cy="150" r="160" fill="url(#sun-glow)" />
          <circle cx="620" cy="150" r="45" fill="#ffffff" />

          {/* 3. Drifting Clouds */}
          <g className="cloud-group-1">
            <path d="M 120 130 Q 135 115 150 130 Q 165 115 180 130 Q 190 140 180 150 L 120 150 Z" fill="#ffffff" opacity="0.85" />
          </g>
          <g className="cloud-group-2">
            <path d="M 450 90 Q 462 78 475 90 Q 488 78 500 90 Q 510 98 500 106 L 450 106 Z" fill="#ffffff" opacity="0.7" />
          </g>

          {/* 4. Far Mountain Range */}
          <polygon points="-50,380 120,200 280,380" fill="url(#mount-grad-far)" opacity="0.8" />
          <polygon points="120,200 145,230 110,230" fill="#ffffff" /> {/* Far mount snow cap */}
          <polygon points="180,380 340,160 520,380" fill="url(#mount-grad-far)" opacity="0.8" />
          <polygon points="340,160 375,205 320,205" fill="#ffffff" /> {/* Far mount snow cap */}

          {/* 5. Mid Mountain Range */}
          <polygon points="60,400 240,240 420,400" fill="url(#mount-grad-mid)" opacity="0.9" />
          <polygon points="240,240 270,280 220,280" fill="#ffffff" /> {/* Mid mount snow cap */}
          <polygon points="320,400 480,220 660,400" fill="url(#mount-grad-mid)" opacity="0.9" />
          <polygon points="480,220 515,265 455,265" fill="#ffffff" /> {/* Mid mount snow cap */}

          {/* 6. Deep Pine Forest Row (Background) */}
          {/* Left row */}
          <polygon points="15,410 25,370 35,410" fill="#074c4f" />
          <polygon points="30,410 40,365 50,410" fill="#074c4f" />
          <polygon points="45,410 58,355 70,410" fill="#074c4f" />
          <polygon points="65,410 75,375 85,410" fill="#074c4f" />
          <polygon points="80,410 92,360 105,410" fill="#074c4f" />
          <polygon points="100,410 110,380 120,410" fill="#074c4f" />
          {/* Right row */}
          <polygon points="660,410 672,360 685,410" fill="#074c4f" />
          <polygon points="680,410 690,375 700,410" fill="#074c4f" />
          <polygon points="695,410 708,355 720,410" fill="#074c4f" />
          <polygon points="715,410 725,365 735,410" fill="#074c4f" />
          <polygon points="730,410 742,370 755,410" fill="#074c4f" />

          {/* 7. Foreground Snowy Slopes & Trails */}
          {/* Dynamic sweeping hills */}
          <path d="M -20 400 Q 200 370 420 410 T 820 380 L 820 520 L -20 520 Z" fill="#f8fafc" />
          <path d="M -20 430 Q 300 420 820 440 L 820 520 L -20 520 Z" fill="#ffffff" />

          {/* Groomed Ski Trail - Corduroy Lines */}
          <path d="M 400 405 L 120 500" stroke="#cbd5e1" strokeWidth="2.5" strokeDasharray="6,4" />
          <path d="M 405 405 L 260 500" stroke="#cbd5e1" strokeWidth="2" strokeDasharray="6,4" />
          <path d="M 410 405 L 400 500" stroke="#cbd5e1" strokeWidth="1.5" strokeDasharray="6,4" />
          <path d="M 415 405 L 540 500" stroke="#cbd5e1" strokeWidth="2" strokeDasharray="6,4" />
          <path d="M 420 405 L 680 500" stroke="#cbd5e1" strokeWidth="2.5" strokeDasharray="6,4" />

          {/* ==========================================================================
             ATHLETE 2: Background Skier (Yellow/Orange Suit)
             ========================================================================== */}
          <g className="skier-athlete background-athlete">
            {/* Shadow under skier */}
            <ellipse cx="525" cy="300" rx="35" ry="5" fill="#cbd5e1" opacity="0.6" />

            {/* Left/Back Ski */}
            <g className="ski-node">
              <polygon points="460,285 530,305 528,308 458,288" fill="#ff9f43" stroke="#ffffff" strokeWidth="1" />
            </g>
            {/* Right/Forward Ski (Splayed Out) */}
            <g className="ski-node">
              <polygon points="490,302 575,292 576,295 491,305" fill="#ff9f43" stroke="#ffffff" strokeWidth="1" />
            </g>

            {/* Left/Back Pole */}
            <g className="pole-node">
              <line x1="508" y1="210" x2="475" y2="295" stroke="#4b5563" strokeWidth="1.5" strokeLinecap="round" />
              {/* Basket */}
              <circle cx="478" cy="287" r="3.5" fill="none" stroke="#ffffff" strokeWidth="1" />
            </g>

            {/* Legs */}
            <path d="M 515 235 L 500 270 L 485 292" fill="none" stroke="#e056fd" strokeWidth="5.5" strokeLinecap="round" strokeLinejoin="round" /> {/* Back leg */}
            <path d="M 525 235 L 535 270 L 545 298" fill="none" stroke="#ff7f50" strokeWidth="6" strokeLinecap="round" strokeLinejoin="round" /> {/* Front leg */}

            {/* Torso & Suit (Yellow/Orange) */}
            <path d="M 512 180 Q 528 175 535 195 L 525 240 L 510 235 Z" fill="#ffa801" stroke="#ffffff" strokeWidth="1.5" />
            <path d="M 514 184 L 521 225" stroke="#ff3f34" strokeWidth="2.5" /> {/* Red speed stripe */}

            {/* Right/Forward Pole */}
            <g className="pole-node">
              <line x1="538" y1="215" x2="568" y2="290" stroke="#4b5563" strokeWidth="1.5" strokeLinecap="round" />
              {/* Basket */}
              <circle cx="565" cy="283" r="3.5" fill="none" stroke="#ffffff" strokeWidth="1" />
            </g>

            {/* Arms */}
            <path d="M 514 185 L 502 205 L 508 215" fill="none" stroke="#ffa801" strokeWidth="4.5" strokeLinecap="round" /> {/* Back arm */}
            <path d="M 528 188 L 542 205 L 538 218" fill="none" stroke="#ffa801" strokeWidth="5.5" strokeLinecap="round" /> {/* Front arm */}

            {/* Head, Beanie & Goggles */}
            <circle cx="522" cy="165" r="11" fill="#f8c291" />
            <path d="M 512 165 C 512 150, 532 150, 532 165 Z" fill="#ff3f34" /> {/* Orange hat */}
            <circle cx="522" cy="148" r="3" fill="#ffffff" /> {/* Pom-pom */}
            <rect x="519" y="125" width="13" height="7" rx="3.5" fill="url(#goggle-back)" stroke="#ffffff" strokeWidth="1.2" transform="rotate(12 522 130)" /> {/* Goggles */}
            {/* Goggle Reflection Glare */}
            <g clipPath="url(#goggle-clip-back)">
              <g className="lens-sheen-group">
                <rect x="510" y="115" width="6" height="30" fill="#ffffff" opacity="0.5" transform="skewX(-20)" />
              </g>
            </g>
          </g>

          {/* ==========================================================================
             ATHLETE 1: Foreground Skier (Navy/Cyan Suit)
             ========================================================================== */}
          <g className="skier-athlete foreground-athlete">
            {/* Shadow under skier */}
            <ellipse cx="325" cy="375" rx="55" ry="7.5" fill="#94a3b8" opacity="0.45" />

            {/* Left/Back Ski */}
            <g className="ski-node">
              <polygon points="215,350 330,380 328,384 213,354" fill="#00f2fe" stroke="#ffffff" strokeWidth="1.5" />
            </g>
            {/* Right/Forward Ski (Splayed Out dynamically) */}
            <g className="ski-node">
              <polygon points="275,378 400,360 401,364 276,382" fill="#00f2fe" stroke="#ffffff" strokeWidth="1.5" />
            </g>

            {/* Left/Back Pole */}
            <g className="pole-node">
              <line x1="302" y1="230" x2="242" y2="368" stroke="#1e293b" strokeWidth="2.5" strokeLinecap="round" />
              {/* Basket */}
              <circle cx="247" cy="356" r="6" fill="none" stroke="#ffffff" strokeWidth="2" />
              <line x1="241" y1="356" x2="253" y2="356" stroke="#ffffff" strokeWidth="1.5" />
            </g>

            {/* Legs */}
            <path d="M 310 270 L 285 320 L 260 355" fill="none" stroke="#0c2461" strokeWidth="10.5" strokeLinecap="round" strokeLinejoin="round" /> {/* Back leg */}
            <path d="M 330 270 L 348 322 L 362 368" fill="none" stroke="#10ac84" strokeWidth="11" strokeLinecap="round" strokeLinejoin="round" /> {/* Front leg */}
            <path d="M 330 270 L 348 322 L 362 368" fill="none" stroke="#0a3d62" strokeWidth="8" strokeLinecap="round" strokeLinejoin="round" /> {/* Front leg inside color */}

            {/* Torso & Suit (Navy with Neon Stripe) */}
            <path d="M 305 185 Q 330 175 342 205 L 328 275 L 305 270 Z" fill="#0a3d62" stroke="#ffffff" strokeWidth="2.2" />
            <path d="M 307 190 Q 322 195 316 250" fill="none" stroke="#00ebfe" strokeWidth="4.5" strokeLinecap="round" /> {/* Cyber Cyan stripe */}

            {/* Right/Forward Pole */}
            <g className="pole-node">
              <line x1="350" y1="238" x2="395" y2="360" stroke="#1e293b" strokeWidth="2.5" strokeLinecap="round" />
              {/* Basket */}
              <circle cx="390" cy="347" r="6" fill="none" stroke="#ffffff" strokeWidth="2" />
              <line x1="384" y1="347" x2="396" y2="347" stroke="#ffffff" strokeWidth="1.5" />
            </g>

            {/* Arms */}
            <path d="M 307 192 L 290 220 L 300 236" fill="none" stroke="#0a3d62" strokeWidth="8" strokeLinecap="round" /> {/* Back arm */}
            <path d="M 332 198 L 354 222 L 348 242" fill="none" stroke="#0a3d62" strokeWidth="9.5" strokeLinecap="round" /> {/* Front arm */}

            {/* Head, Headband & Hair */}
            <circle cx="321" cy="162" r="16" fill="#f8c291" stroke="#ffffff" strokeWidth="1.5" />
            <path d="M 307 160 C 307 140, 335 140, 335 160 Z" fill="#ff7f50" /> {/* Orange headband */}
            <path d="M 307 165 C 302 165, 296 175, 298 185" fill="none" stroke="#ff7f50" strokeWidth="4.5" strokeLinecap="round" /> {/* Hair whip */}
            
            {/* Cyber Goggles */}
            <rect x="317" y="145" width="22" height="12" rx="6" fill="url(#goggle-fore)" stroke="#ffffff" strokeWidth="2" transform="rotate(15 321 162)" />
            {/* Goggle Reflection Glare */}
            <g clipPath="url(#goggle-clip-fore)">
              <g className="lens-sheen-group">
                <rect x="300" y="130" width="10" height="40" fill="#ffffff" opacity="0.65" transform="skewX(-20)" />
              </g>
            </g>
          </g>
        </svg>
      </div>

      {/* Modern Poster Style Solid Bottom Footer Banner */}
      <footer className="poster-footer">
        <div className="footer-content">
          <div className="text-banner">
            <h1 className="footer-title">404 - lost in the powder</h1>
            <p className="footer-subtitle">cross-country skiing / sci di fondo / no trail found</p>
          </div>
          
          <div className="button-banner">
            <Link to="/" className="trailhead-button">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" className="home-icon">
                <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
                <polyline points="9 22 9 12 15 12 15 22" />
              </svg>
              <span>Back to Trailhead</span>
            </Link>
          </div>
        </div>
      </footer>

      {/* Atmospheric Winter Snowflake Layer */}
      <div className="snowflakes-container">
        {snowflakes}
      </div>
    </div>
  );
};

export default NotFound;
