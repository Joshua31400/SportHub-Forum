document.addEventListener('DOMContentLoaded', function() {
    const fileInput = document.getElementById('image');
    const fileChosen = document.getElementById('file-chosen');
    const fileUploadBtn = document.querySelector('.file-upload-btn');
    const fileContainer = document.querySelector('.custom-file-upload');

    const style = document.createElement('style');
    style.textContent = `
        .file-success { color: #28a745 !important; font-weight: 500; }
        .file-error { color: #dc3545 !important; font-weight: 500; }
        .btn-success { background-color: #28a745 !important; border-color: #28a745 !important; }
        .btn-error { background-color: #dc3545 !important; border-color: #dc3545 !important; }
        .container-success { 
            border: 1px solid #28a745; 
            border-radius: 5px; 
            padding: 8px; 
            background-color: rgba(40, 167, 69, 0.1); 
        }
        .container-error { 
            border: 1px solid #dc3545; 
            border-radius: 5px; 
            padding: 8px; 
            background-color: rgba(220, 53, 69, 0.1); 
        }
    `;
    document.head.appendChild(style);

    if (fileInput && fileChosen && fileUploadBtn) {
        updateFileStatus();

        fileInput.addEventListener('change', updateFileStatus);

        function updateFileStatus() {
            if (fileInput.files && fileInput.files[0]) {
                fileChosen.textContent = fileInput.files[0].name;
                fileChosen.className = 'file-success';
                fileUploadBtn.className = 'file-upload-btn btn-success';
                fileContainer.className = 'custom-file-upload container-success';
            } else {
                fileChosen.textContent = 'No file selected';
                fileChosen.className = 'file-error';
                fileUploadBtn.className = 'file-upload-btn btn-error';
                fileContainer.className = 'custom-file-upload container-error';
            }
        }
    }
});