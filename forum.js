document.addEventListener('DOMContentLoaded', function() {
    const newTopicBtn = document.querySelector('.new-topic-btn');
    const modal = document.getElementById('new-topic-modal');
    const closeModal = document.querySelector('.close-modal');

    newTopicBtn.addEventListener('click', function() {
        modal.style.display = 'block';
    });

    closeModal.addEventListener('click', function() {
        modal.style.display = 'none';
    });

    window.addEventListener('click', function(event) {
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });

    const newTopicForm = document.getElementById('new-topic-form');
    newTopicForm.addEventListener('submit', function(event) {
        event.preventDefault();
        const topicTitle = document.getElementById('topic-title').value;
        const topicContent = document.getElementById('topic-content').value;

        if (topicTitle && topicContent) {
            console.log('New Topic Created:', { title: topicTitle, content: topicContent });
            modal.style.display = 'none';
            newTopicForm.reset();
            alert('New topic created successfully!');
        } else {
            alert('Please fill in all fields.');
        }
    });
});