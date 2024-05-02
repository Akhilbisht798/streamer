## Streamer
Stream Directly from Browser.

### Functinal Requirement
- Stream your video and screen directly from browser.
- Can control video, audio and screen.

### System Design
- ALB load balancer in front of service.
- ECS service with 2 instance.
- maybe a auto scaling group to scale with every user.
