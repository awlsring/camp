use std::{marker::PhantomData, task::{Context, Poll}, sync::Arc};

use aws_smithy_http_server::{body::BoxBody, response::{Response, IntoResponse}};
use http_body::{Body as HttpBody};
use tower::Service;

use super::{exception::UnauthorizedException, controller::AuthController};

pub struct AuthService<Protocol, Operation, Service, Config> {
    inner: Service,
    controller: Arc<AuthController>,
    operation: &'static str,
    config: Config,
    _protocol: PhantomData<Protocol>,
    _operation: PhantomData<Operation>,
}

impl <Protocol, Operation, Service, Config> AuthService<Protocol, Operation, Service, Config> {
    pub fn new(inner: Service, operation: &'static str, auth: Arc<AuthController>, config: Config) -> AuthService<Protocol, Operation, Service, Config> {
        AuthService {
            inner,
            controller: auth,
            operation,
            config,
            _protocol: PhantomData,
            _operation: PhantomData,
        }
    }
}

impl <Protocol, Operation, Service, Config> Clone for AuthService<Protocol, Operation, Service, Config>
where
    Service: Clone,
    Config: Clone,
{
    fn clone(&self) -> Self {
        AuthService {
            inner: self.inner.clone(),
            controller: self.controller.clone(),
            operation: self.operation,
            config: self.config.clone(),
            _protocol: PhantomData,
            _operation: PhantomData,
        }
    }
}

impl <Body, Protocol, Operation, S, Config> Service<http::Request<Body>> for AuthService<Protocol, Operation, S, Config>
where
    Body: HttpBody + Send + Sync + 'static,
    Body::Data: Send,

    S: Service<http::Request<Body>, Response = Response<BoxBody>> + Send,
    S::Future: Send,
    S: Clone + 'static,

    UnauthorizedException: IntoResponse<Protocol>,
{
    type Response = S::Response;
    type Error = S::Error;
    type Future = std::pin::Pin<
        Box<dyn std::future::Future<Output = Result<Self::Response, Self::Error>> + Send>
    >;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }

    fn call(&mut self, req: http::Request<Body>) -> Self::Future {

        let inner = self.inner.clone();
        let mut inner = std::mem::replace(&mut self.inner, inner);

        let auth_controller = self.controller.clone();
        let op = self.operation;

        let f = async move {
            let headers = req.headers();
            let auth_header = headers.get("authorization");

            if auth_controller.auth(op, auth_header).await {
                return inner.call(req).await;
            }
            
            // This sucks, fix later
            Ok(IntoResponse::<Protocol>::into_response(UnauthorizedException))
        };
        
        Box::pin(f)
    }
}