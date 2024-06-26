## Streamer
Stream Directly from Browser.

### Functinal Requirement
- Stream your video and screen directly from browser.
- Can control video, audio and screen.

### System Design
- ALB load balancer in front of service.
- ECS service with 2 instance.
- maybe a auto scaling group to scale with every user.


## Help
- ECS has two roles: Task Role and Task Execution Role. The Task Role is the perms the application uses. The Task Execution Role is the role the ECS servce assumes to pull containers and write logs.
- Go to the Task Definition, and check to see if the Task Execution Role has a trust relationship to the principal: ecs-tasks.amazonaws.com and has the AmazonECSTaskExecutionRolePolicy managed policy applied.

- [youtube](https://www.youtube.com/watch?v=fb2zJlcE1bE)
