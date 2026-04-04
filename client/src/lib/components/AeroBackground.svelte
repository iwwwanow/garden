<div class="bg">
  <div class="sky" />
  <div class="blob b1" />
  <div class="blob b2" />
  <div class="blob b3" />
  <div class="blob b4" />
  <div class="blob b5" />
  <div class="blob b6" />
  <div class="flare" />
</div>

<style>
.bg {
  position: fixed;
  inset: 0;
  z-index: -1;
  overflow: hidden;
  pointer-events: none;
}

.sky {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    180deg,
    #7ec8e3 0%,
    #b8dff0 28%,
    #cceef4 55%,
    #d8f4ea 80%,
    #e2f8f0 100%
  );
}

/* Bokeh blobs */
.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(70px);
  animation: drift linear infinite;
  will-change: transform;
}

.b1 {
  width: 380px; height: 380px;
  background: radial-gradient(circle, rgba(64,244,193,0.55) 0%, transparent 70%);
  top: -120px; left: -100px;
  animation-duration: 22s;
}
.b2 {
  width: 300px; height: 300px;
  background: radial-gradient(circle, rgba(64,154,244,0.45) 0%, transparent 70%);
  top: 60px; right: -80px;
  animation-duration: 29s;
  animation-delay: -9s;
}
.b3 {
  width: 240px; height: 240px;
  background: radial-gradient(circle, rgba(255,255,255,0.8) 0%, transparent 70%);
  top: 38%; left: 25%;
  animation-duration: 35s;
  animation-delay: -17s;
}
.b4 {
  width: 420px; height: 420px;
  background: radial-gradient(circle, rgba(184,224,255,0.5) 0%, transparent 70%);
  bottom: 60px; left: -130px;
  animation-duration: 27s;
  animation-delay: -6s;
}
.b5 {
  width: 200px; height: 200px;
  background: radial-gradient(circle, rgba(168,240,216,0.6) 0%, transparent 70%);
  bottom: 180px; right: 10px;
  animation-duration: 24s;
  animation-delay: -13s;
}
.b6 {
  width: 160px; height: 160px;
  background: radial-gradient(circle, rgba(255,220,180,0.35) 0%, transparent 70%);
  top: 15%; left: 55%;
  animation-duration: 40s;
  animation-delay: -22s;
}

/* Lens flare */
.flare {
  position: absolute;
  width: 140px; height: 140px;
  top: 3%; right: 12%;
  background: radial-gradient(
    circle,
    rgba(255,255,255,0.95) 0%,
    rgba(255,240,160,0.5) 35%,
    transparent 70%
  );
  filter: blur(6px);
  animation: flare 10s ease-in-out infinite;
}

@keyframes drift {
  0%   { transform: translate(0, 0)     scale(1);    }
  20%  { transform: translate(18px, -22px) scale(1.06); }
  45%  { transform: translate(-12px, 14px) scale(0.94); }
  70%  { transform: translate(-22px, -18px) scale(1.03); }
  100% { transform: translate(0, 0)     scale(1);    }
}

@keyframes flare {
  0%, 100% { opacity: 0.65; transform: scale(1); }
  50%       { opacity: 1;    transform: scale(1.15); }
}
</style>
