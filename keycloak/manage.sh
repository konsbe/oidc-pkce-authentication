#!/bin/bash

# Simple management script for VOMT Keycloak Local Deployment

case "$1" in
    start)
        echo "🚀 Starting VOMT Keycloak Local Deployment..."
        docker compose up -d
        echo "⏳ Waiting for Keycloak to be ready..."
        sleep 60
        echo "⚙️  Configuring Keycloak with VOMT setup..."
        bash configure-keycloak.sh
        ;;
    stop)
        echo "🛑 Stopping VOMT Keycloak..."
        docker compose down
        ;;
    restart)
        echo "🔄 Restarting VOMT Keycloak..."
        docker compose down
        docker compose up -d
        sleep 60
        bash configure-keycloak.sh
        ;;
    status)
        echo "📊 Container Status:"
        docker compose ps
        ;;
    logs)
        docker compose logs -f
        ;;
    clean)
        echo "🧹 Cleaning up (this will remove all data)..."
        read -p "Are you sure? This will delete all Keycloak data [y/N]: " confirm
        if [[ $confirm =~ ^[Yy]$ ]]; then
            docker compose down -v
            docker system prune -f
            echo "✅ Cleanup completed."
        else
            echo "❌ Cleanup cancelled."
        fi
        ;;
    config)
        echo "⚙️  Running Keycloak configuration only..."
        bash configure-keycloak.sh
        ;;
    test)
        echo "🧪 Testing authentication..."
        curl -s -X POST "http://localhost:8080/access/realms/vomt/protocol/openid-connect/token" \
          -H "Content-Type: application/x-www-form-urlencoded" \
          -d "username=vomtadmin" \
          -d "password=Admin@123" \
          -d "grant_type=password" \
          -d "client_id=spog" | grep -o '"access_token":"[^"]*"' | head -1
        ;;
    *)
        echo "🐳 VOMT Keycloak Local Deployment Manager"
        echo ""
        echo "Usage: $0 {start|stop|restart|status|logs|clean|config|test}"
        echo ""
        echo "Commands:"
        echo "  start   - Start Keycloak and configure VOMT realm"
        echo "  stop    - Stop Keycloak"
        echo "  restart - Restart Keycloak and reconfigure"
        echo "  status  - Show container status"
        echo "  logs    - Show container logs"
        echo "  clean   - Remove all containers and data"
        echo "  config  - Run configuration script only"
        echo "  test    - Test authentication with vomtadmin"
        echo ""
        echo "Access URLs:"
        echo "  🌐 Keycloak: http://localhost:8080/access"
        echo "  🔧 Admin: http://localhost:8080/access/admin (admin/admin)"
        echo "  🏠 VOMT Realm: http://localhost:8080/access/realms/vomt"
        exit 1
        ;;
esac
