(function () {
  const root = document.documentElement;
  let rafId = 0;

  function animateGlow() {
    const time = Date.now() * 0.00012;
    const x = Math.sin(time) * 12;
    const y = Math.cos(time * 1.3) * 10;
    root.style.setProperty("--glow-shift-x", `${x}px`);
    root.style.setProperty("--glow-shift-y", `${y}px`);
    rafId = window.requestAnimationFrame(animateGlow);
  }

  animateGlow();

  window.addEventListener("beforeunload", function () {
    if (rafId) {
      window.cancelAnimationFrame(rafId);
    }
  });
})();
