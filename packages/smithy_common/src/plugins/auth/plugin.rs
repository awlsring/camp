/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

//! Provides an example [`Plugin`] implementation - [`PrintPlugin`].

use std::sync::Arc;

use aws_smithy_http_server::{
    operation::{OperationShape},
    plugin::{Plugin, PluginPipeline, PluginStack},
};

use super::{controller::AuthController, service::AuthService};

#[derive(Clone, Debug)]
pub struct AuthPlugin<Config> {
    controller: Arc<AuthController>,
    config: Config,
}

impl <Config> AuthPlugin<Config> {
    pub fn new(controller: Arc<AuthController>, config: Config) -> Self {
        AuthPlugin { controller, config }
    }
}

impl <Protocol, Op, Svc, Config> Plugin<Protocol, Op, Svc> for AuthPlugin<Config>
where
Op: OperationShape,
    Config: Clone,
{
    type Service = AuthService<Protocol, Op, Svc, Config>;

    fn apply(&self, inner: Svc) -> Self::Service {
        AuthService::new(inner, Op::ID.name(), self.controller.clone(), self.config.clone())
    }
}

pub trait AuthExtension<ExistingPlugins> {
    fn auth<Config>(self, controller: Arc<AuthController>, config: Config) -> PluginPipeline<PluginStack<AuthPlugin<Config>, ExistingPlugins>>;
}

impl <ExistingPlugins> AuthExtension<ExistingPlugins> for PluginPipeline<ExistingPlugins> {
    fn auth<Config>(self, controller: Arc<AuthController>, config: Config) -> PluginPipeline<PluginStack<AuthPlugin<Config>, ExistingPlugins>> {
        self.push(AuthPlugin::new(controller, config))
    }
}