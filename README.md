# RabbitMQ Image Processor

RabbitMQ Image processor is a homepage and a processor entity that generates image
previews of uploaded images. The image previews are generated using graphicsmagick.
The process load is distributed using rabbitMQ and the processor services communicates with
rabbit mq for images. The image preview names are stored as the adler32 checksum,
so that a database doesn't have to store information about every single file.

## Image processing

Images are only processed into a single image format. They are scaled down to
size such that the maximum width and maximum height is 700 px.