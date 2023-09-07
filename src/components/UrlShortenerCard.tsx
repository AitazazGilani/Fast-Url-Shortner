import React, { useState } from 'react';
import { Container, Card, CardBody, Heading, Input, Button } from '@shadcn/ui';

function UrlShortenerCard() {
  const [originalUrl, setOriginalUrl] = useState('');
  const [shortenedUrl, setShortenedUrl] = useState('');

  const handleShortenClick = async () => {
    // Implement your URL shortening logic here, e.g., make an API call.

    // For this example, we're setting a placeholder shortened URL.
    setShortenedUrl('http://shortened-link.com/abcdef');
  };

  return (
    <Container className="flex items-center justify-center h-screen">
      <Card>
        <CardBody>
          <Heading as="h2" size="lg" className="mb-4">
            URL Shortener
          </Heading>
          <Input
            type="text"
            placeholder="Enter your URL"
            value={originalUrl}
            onChange={(e) => setOriginalUrl(e.target.value)}
            className="mb-4"
          />
          <Button onClick={handleShortenClick} className="mb-4">
            Shorten URL
          </Button>
          {shortenedUrl && (
            <p>
              Shortened URL: <a href={shortenedUrl}>{shortenedUrl}</a>
            </p>
          )}
        </CardBody>
      </Card>
    </Container>
  );
}

export default UrlShortenerCard;
