function openModal() {
    document.getElementById("modal").style.display = "block";
}
  
// Close the Modal
function closeModal() {
    document.getElementById("modal").style.display = "none";
}

var slideIndex = 0;
showSlides(slideIndex);
  
// Next/previous controls
function plusSlides(n) {
    showSlides(slideIndex += n);
}
  
// Thumbnail image controls
function currentSlide(n) {
    showSlides(slideIndex = n);
}
  
function showSlides(n) {
    var i;
    var slides = document.getElementsByClassName("slides");
    if (n > slides.length-1) {slideIndex = 0}
    if (n < 0) {slideIndex = slides.length-1}
    for (i = 0; i < slides.length; i++) {
      slides[i].style.display = "none";
    }
    slides[slideIndex].style.display = "block";
}