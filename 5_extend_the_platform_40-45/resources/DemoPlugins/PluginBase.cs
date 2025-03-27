using System;
using System.ServiceModel;
using Microsoft.Xrm.Sdk;

namespace DemoPlugins
{
    public abstract class PluginBase : IPlugin
    {
        protected IPluginExecutionContext context;
        protected IOrganizationService organisationService;
        protected ITracingService tracingService;

        /// <summary>
        /// Contains the business logic to be executed by the plugin. This
        /// method must be implemented by derived classes to define the specific
        /// business logic for the plugin.
        /// </summary>
        /// <param name="serviceProvider">
        /// The service provider that provides access to the services required by the plugin.
        /// </param>
        public abstract void ExecuteBusinessLogic(IServiceProvider serviceProvider);

        /// <summary>
        /// Executes the plugin. This method is called by the Dynamics 365
        /// framework when the plugin is triggered.
        /// </summary>
        /// <param name="serviceProvider">
        /// The service provider that provides access to the services required
        /// by the plugin.
        /// </param>
        /// <remarks>
        /// This method wraps the execution of the business logic in a try-catch
        /// block to handle any exceptions that occur during the execution. If
        /// an exception is thrown, it is caught and rethrown as an
        /// InvalidPluginExecutionException.
        /// </remarks>
        public void Execute(IServiceProvider serviceProvider)
        {
            try
            {
                context = TryGetService<IPluginExecutionContext>(serviceProvider);
                var orgSvcFactory = TryGetService<IOrganizationServiceFactory>(serviceProvider);
                organisationService = orgSvcFactory.CreateOrganizationService(context.UserId);
                tracingService = TryGetService<ITracingService>(serviceProvider);
                ExecuteBusinessLogic(serviceProvider);
            }
            catch (Exception ex)
            {
                throw new InvalidPluginExecutionException(ex.Message, ex);
            }
        }

        /// <summary>
        /// Ensures that the provided service provider is not null.
        /// </summary>
        /// <param name="serviceProvider">
        /// The service provider to check.
        /// </param>
        /// <exception cref="ArgumentNullException">
        /// Thrown when the service provider is null.
        /// </exception>
        private void GuardServiceProviderNotNull(IServiceProvider serviceProvider)
        {
            if (serviceProvider is null)
            {
                throw new ArgumentNullException("Service provider is null");
            }
        }

        /// <summary>
        /// Attempts to retrieve a service of the specified type from the
        /// service provider.
        /// </summary>
        /// <typeparam name="T">
        /// The type of service to retrieve.
        /// </typeparam>
        /// <param name="serviceProvider">
        /// The service provider that provides access to the services.
        /// </param>
        /// <returns>
        /// The service of the specified type if found; otherwise, throws an
        /// exception.
        /// </returns>
        /// <exception cref="ArgumentNullException">
        /// Thrown when the service provider is null.
        /// </exception>
        /// <exception cref="ArgumentException">
        /// Thrown when the service of the specified type cannot be accessed
        /// from the service provider.
        /// </exception>
        protected T TryGetService<T>(IServiceProvider serviceProvider)
            where T : class
        {
            GuardServiceProviderNotNull(serviceProvider);
            if (serviceProvider.GetService(typeof(T)) is T service)
            {
                return service;
            }
            throw new ArgumentException(
                $"Unable to access ${typeof(T).Name} from the service provider"
            );
        }

        /// <summary>
        /// Retrieves a parameter from the context parameter collection or
        /// returns the default value if the parameter is not found.
        /// </summary>
        /// <typeparam name="T">
        /// The type of the parameter to retrieve.
        /// </typeparam>
        /// <param name="paramCollection">
        /// The collection of parameters.
        /// </param>
        /// <param name="key">
        /// The key of the parameter to retrieve.
        /// </param>
        /// <returns>
        /// The parameter of the specified type if found; otherwise, the default
        /// value of the specified type.
        /// </returns>
        protected T GetContextParameterOrDefault<T>(ParameterCollection paramCollection, string key)
        {
            if (paramCollection.TryGetValue(key, out T value))
            {
                return value;
            }
            return default;
        }

        /// <summary>
        /// Retrieves an entity image from the specified image collection by its
        /// name.
        /// </summary>
        /// <param name="imageCollection">
        /// The collection of entity images.
        /// </param>
        /// <param name="imageName">
        /// The name of the entity image to retrieve.
        /// </param>
        /// <returns>
        /// The entity image if found; otherwise, <c>null</c>.
        /// </returns>
        protected Entity GetEntityImage(EntityImageCollection imageCollection, string imageName)
        {
            if (imageCollection.TryGetValue(imageName, out Entity e))
            {
                return e;
            }
            return null;
        }
    }
}
