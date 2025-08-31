import { useState } from 'react';
import axios from 'axios';

export function Upload () {
  const [file, setFile] = useState<File | null>(null);
  const [progress, setProgress] = useState('');

  const handleUpload = async () => {
    const formData = new FormData();
    formData.append('file', file!);
    // TODO Use react-query.
    const res = await axios.post('/api/upload', formData);
    const id = res.data.id;

    const es = new EventSource(`/api/stream/${id}`);
    es.onmessage = (e) => setProgress(e.data);
  };

  return (
    <div>
      <input type="file" onChange={(e) => setFile(e.target.files?.[0] as File)} />
      <button onClick={handleUpload}>Upload</button>
      <p>Progress: {progress}</p>
    </div>
  );
};

